# treetop

<img src="treetop.png" width="400">

![build](https://github.com/eigenhombre/treetop/actions/workflows/build.yml/badge.svg)

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


<img src="treetop.gif">

# License

Copyright © 2022, John Jacobsen. MIT License.

# Disclaimer

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
