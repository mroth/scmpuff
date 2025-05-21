IMAGE = scmpuff-demo-recorder

.PHONY: image
image:
	docker build -t $(IMAGE) .

demo-raw.gif: demo.tape
	docker run --rm -v $(PWD):/vhs $(IMAGE) $<

demo.gif: demo-raw.gif
	gifsicle -O3 $< -o $@

.PHONY: demo
demo: demo.gif

.PHONY: demo-web
demo-web:
	docker run --rm -v $(PWD):/vhs $(IMAGE) demo.web.tape

.PHONY: clean
clean:
	rm -f demo-raw.gif
