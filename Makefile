deploy:
	gcloud functions deploy GOLFS \
		--runtime go113 \
		--trigger-http \
		--env-vars-file .env.yaml
