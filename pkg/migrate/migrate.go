package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/katallaxie/autobot/pkg/migrate/cql"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	cm "github.com/scylladb/gocqlx/v2/migrate"
)

type migrate struct {
	keyspace string
}

// New ...
func New(keyspace string) *migrate {
	m := new(migrate)
	m.keyspace = keyspace

	return m
}

// Migrate ...
func (m *migrate) Migrate(ctx context.Context, cluster *gocql.ClusterConfig) error {
	err := m.createKeySpace(cluster, m.keyspace) //nolint:contextcheck
	if err != nil {
		return err
	}
	cluster.Keyspace = m.keyspace

	session, err := gocqlx.WrapSession(cluster.CreateSession()) //nolint:contextcheck
	if err != nil {
		return err
	}
	defer session.Close()

	log := func(ctx context.Context, session gocqlx.Session, ev cm.CallbackEvent, name string) error {
		return nil
	}

	reg := cm.CallbackRegister{}
	reg.Add(cm.BeforeMigration, "m1.cql", log)
	reg.Add(cm.AfterMigration, "m1.cql", log)
	reg.Add(cm.CallComment, "1", log)
	cm.Callback = reg.Callback

	if err := cm.FromFS(ctx, session, cql.Files); err != nil {
		return err
	}

	return nil
}

func (m *migrate) createKeySpace(cluster *gocql.ClusterConfig, keyspace string) error {
	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second

	session, err := gocqlx.WrapSession(c.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()

	{
		err := session.ExecStmt(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : %d}`, keyspace, 1))
		if err != nil {
			return fmt.Errorf("create keyspace: %w", err)
		}
	}

	return nil
}
