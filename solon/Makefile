copy-certs-locally:
	mkdir -p /tmp/config/vault/
	kubectl -n odysseia get secret vault-server-tls -o jsonpath='{.data.vault\.key}' | base64 --decode > /tmp/config/vault/vault.key
	kubectl -n odysseia get secret vault-server-tls -o jsonpath='{.data.vault\.crt}' | base64 --decode > /tmp/config/vault/vault.crt
	kubectl -n odysseia get secret vault-server-tls -o jsonpath='{.data.vault\.ca}' | base64 --decode > /tmp/config/vault/vault.ca

copy-solon-certs-locally:
	cd ../ && cp eratosthenes/fixture/solon/tls.crt /tmp/tls.crt && cp eratosthenes/fixture/solon/tls.key /tmp/tls.key