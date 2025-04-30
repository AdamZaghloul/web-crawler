package main

import "testing"

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		inputBody string
		expected string
	}{
		{
			name:     "simple test",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html>
					<body>
						<a href="https://blog.boot.dev">
						<span>Go to Boot.dev, you React Andy</span>
						</a>
					</body>
					</html>`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name:     "two links one relative",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html>
							<body>
								<a href="/path/one">
									<span>Boot.dev</span>
								</a>
								<a href="https://other.com/path/one">
									<span>Boot.dev</span>
								</a>
							</body>
						</html>`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "three links one relative",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html>
							<body>
								<a href="/path/one">
									<span>Boot.dev</span>
								</a>
								<a href="https://other.com/path/one">
									<span>Boot.dev</span>
								</a>
								<div>
									<a href="https://other.com/path/two">
										<span>Boot.dev</span>
									</a>
								</div>
							</body>
						</html>`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one", "https://other.com/path/two"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.input)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
