# treetop

<img src="treetop.png" width="400">

Treetop is a fast, animated "disk summarizer."

Essentially, a replacement for

    du -ks * | sort -rn

... but animated, so you can get progressive
feedback while your disk is being scanned
(helpful for very large directories).

# Install

Install [Go](https://go.dev/), then

    go get github.com/eigenhombre/treetop

# Examples

    treetop    # Shows current directory usage
    treetop ~  # Shows home directory usage



https://user-images.githubusercontent.com/382668/151611083-f659fde6-3b09-47d4-8272-bcf7909fecc7.mov

