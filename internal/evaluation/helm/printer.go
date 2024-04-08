package helm

import "fmt"

func Print(render *Render) {

	fmt.Println("------------ values files ---------------")
	for _, valueFile := range render.ValueFiles {
		fmt.Println(fmt.Sprintf("file : %s ", valueFile.Name))
		println(valueFile.Data)
	}

	if render.Values != nil {
		fmt.Println("-------------- values -----------------")
		println(render.Values.Data)
	}

	if render.MergedValues != nil {
		fmt.Println("------------ merged values ---------------")
		println(render.MergedValues.Data)
	}

	fmt.Println()
	fmt.Println("------------ generated manifests ---------------")
	if len(render.Manifests) == 0 {
		println("No generated manifest")
	} else {
		PrintManifests(render.Manifests)
	}
}

func PrintManifests(manifests []*Manifest) {

	for _, m := range manifests {

		fmt.Println(m.Content)
		fmt.Println("------")
	}
}
