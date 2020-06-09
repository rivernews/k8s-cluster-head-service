package utilities

func SendSlackMessage(message string) {
	Fetch(FetchOption{
		URL: SendSlackURL,
		Method: "POST",
		PostData: map[string]string{
			"text": message,
			"channel": "#build",
			"username": "Kubernetes Cluster Head Service",
			"icon_url": "https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png",
		},
	})
}
