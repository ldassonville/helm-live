package pseudoyaml

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLookupFirstMatch(t *testing.T) {

	tests := map[string]struct {
		value    string
		path     string
		expected string
	}{
		`test case 1`: {
			path:     "metadata.name",
			expected: "nginx",
			value: `
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
wrong yaml syntax at tail
`,
		},
		`test case 2`: {
			path:     "apiVersion",
			expected: "v1",
			value: `
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
wrong yaml syntax at tail
`,
		},
		`test case 3`: {
			path:     "metadata.name",
			expected: "my-component",
			value: `# Source: helm-expo/templates/cloudflare-service.yaml
apiVersion: my.company/v1alpha1
kind: PodLike
metadata:
  name: my-component
  labels:    
    demo: "toto"
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
`,
		},
	}

	// test attribute declaration
	// here
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			firstMatch := LookupFirstMatch(strings.Split(test.path, "."), test.value)
			assert.Equal(t, test.expected, firstMatch)
		})
	}

}
