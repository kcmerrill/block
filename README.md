![Block](block.png "Block")

# Block

A really simplistic approach to a CLI launcher. Using a very *cough* sophisticated approach to determining which file/folder/script to run. Many uses, but here are a few:

* Open applications via the cli
* Change directories quickly
* Open folders/code with your favorite editor quickly
* Pause/Play spotify via the cli with just a few characters

# Usage

In order to get the full functionality of `block`, you'll need to setup a simple block function in your bash profile. Before doing so, it's recommended that you try out block a few tries to make sure it fits your needs, and to ensure that the logic used to score the inventory is to your liking. 

```sh
echo 'b () { $(block "$@" | tail -1) ; }' >> ~/.bash_profile
```

### Scripts

Any matches that are found within your `$HOME/block/` directory will be run with `bash -c $HOME/lock/matched.file`. This is incredibly useful for applescripts to pause/play spotify for example(see the `scripts` folder), compose email messages, or any other scripted tasks you can think of. I've included my applescript file. When placed in `$HOME/block/spotify` for me, I can simply run `$ b spot` and that will pause/play my spotify when I'm on the CLI which I am most of the day without needing to do other key commands(my work keyboard doesn't have pause/play buttons. :shame:)

```sh
// music playing
$ b spot
// music pauses
```

### Defaults

By default(outside of the scripts directory outlined above) block will attempt to `open` or `cd` depending on the result set it matched on. So if it scored highest a directory, block will change to that directory, or if it's a file type, it will attempt to `open` it. You can however influence it's decesion by passing in the first param of the search to be the application you want to use. A few examples:

```sh
// If you know you are looking for a super cool directory
$ b cd somesupercoolfuzzymatcheddirectory
// switched to that directory if succesful

// Lets open a known file in vim for easy editing
$ b vim inventory.go
// matches on inventory.go and opens it in vim

// unsure? lazy? block will do what it thinks is right
$ b inventory.go // will `open path/to/inventory.go`
```

## Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/block/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/block/linux/amd64)

via go:

`$ go get -u github.com/kcmerrill/block`

## Configuration

In your home directory, you can place a `block/block.yml` file where you can override and boost certain folders or files with key/values. The key is simply something that the file will start with, and the value is the command to override with. Of course, if you specify an action on the CLI that will be used instead of this file. 

Boosting is the same. If the file starts with the key, the value will be boosted by the int value. Boosting will bump up the file in priority.

```yaml
overrides:
  "/code/": "vim"
boost:
  "/code/": 100
```
