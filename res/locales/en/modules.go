package en

import (
	"github.com/katallaxie/autobot/pkg/modules"
)

func init() {
	modules.RegisterModules("en", []modules.Module{
		{
			Tag: modules.K8sTag,
			Patterns: []string{
				"What is the status of cluster ",
				"Ping ",
				"Can you reach the cluster ",
			},
			Responses: []string{
				"%s seems to be up",
				"Yep, %s is still going",
			},
			Replacer: modules.K8sReplacer,
		},

		{
			Tag: modules.MathTag,
			Patterns: []string{
				"Give me the result of ",
				"Calculate ",
			},
			Responses: []string{
				"The result is %s",
				"That makes %s",
			},
			Replacer: modules.MathReplacer,
		},

		{
			Tag: modules.RandomTag,
			Patterns: []string{
				"Give me a random number",
				"Generate a random number",
			},
			Responses: []string{
				"The number is %s",
			},
			Replacer: modules.RandomNumberReplacer,
		},

		{
			Tag: modules.JokesTag,
			Patterns: []string{
				"Tell me a joke",
				"Make me laugh",
			},
			Responses: []string{
				"Here you go, %s",
				"Here's one, %s",
			},
			Replacer: modules.JokesReplacer,
		},
		{
			Tag: modules.AdvicesTag,
			Patterns: []string{
				"Give me an advice",
				"Advise me",
			},
			Responses: []string{
				"Here you go, %s",
				"Here's one, %s",
				"Listen closely, %s",
			},
			Replacer: modules.AdvicesReplacer,
		},
	})
}
