_NAME=go-redir-yourself

compile: 
	for os in darwin linux; do \
		for arch in amd64 arm64; do \
			GOOS=$$os GOARCH=$$arch go build -o bin/$(_NAME) . ; \
			cd bin ; \
			tar czf $(_NAME)-$$os-$$arch.tar.gz $(_NAME) ; \
			sha256sum $(_NAME)-$$os-$$arch.tar.gz > $(_NAME)-$$os-$$arch.tar.gz.sha256 ; \
			cd - ; \
		done ; \
	done

clean:
	rm bin/$(_NAME)*
