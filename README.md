# Vader

**VaDeR: Venvs Done Right**

So we all like Python.  And PIP is pretty cool as a package manager.  But
something that's always pissed me off about Python's ecosystem is that PIP likes
to shoot itself in the foot.  Let's say program A depends on version 1.2 of
package B.  So we install package B version 1.2.  But then we need to go run
program C, and that requires version 1.6 of package B.  But 1.6 isn't perfectly
backwards-compatible with version 1.2, and vice-versa, and things break.

Eventually someone figure out that the best idea to deal with dependency hell
is to make this cool thing called Virtual Envrionments.  Most of the time our
IDEs and other tools take care of making the call to
`source some/path/or/whatever/virtualenv` to set up the evironmental variables
for the current shell so that PIP knows where to install the packages, but this
is messy and annoying and you have to do it for each shell and if you forget
then you have a good chance of breaking things.  It feels like a
horribly-design bandaid to me and I know there's a better way of not shooting
ourselves in the foot with configurations and environmental variables.

We have this nice envvar called `PYTHONPATH` for a reason, but it seems like
nobody uses it outside of special situations.

## How does Vader work?

Instead of doing whatever PIP feels like doing for installing packages, we wrap
the invocation of Python and modify the `PYTHONPATH` environmental variable so
that `python` knows where to find libraries instead of keeping them in just the
same place as the Python stdlib.

You can normally invoke your Python programs using `vader run ./foo.py` and it
should sort everything out on its own.

### Configuration

Vader is capable of looking for `requirements.txt` files and automatically
downloading them into your user's `~/.vader` repository.  *These* dependencies
are a lot smarter and know not to step on eachother's toes.  We collect all of
them together for to figure what to set `PYTHONPATH` to so that `python` knows
where to look for modules.

There's also a `Vaderfile` you can specify using TOML for richer configuration,
which will let you use `vader run` directly without specifying an entry file.

## Usage

* `vader download` :: Downloads dependencies so that you can run the program
without downloading anything.

* `vader run [<entry>]` :: Downloads dependencies and runs codegen as necessary,
and runs the program.  Don't need to specify an entry point if you have a proper
Vaderfile.

* `vader codegen` :: Runs code generation programs to build transient
dependencies.

* `vader clean` :: Cleans transient dependencies direcory.

* `vader package` :: Packages program, all dependencies into a ZIP file, adding
a bootstrap script to run it when executed.

## // TODO

More important at the top.

- [ ] Python binary execution
- [ ] PIP dependency identification, downloading (`requirements.txt`)
- [ ] Figuring out how to use `Vaderfile`s
- [ ] Packaging modular/monolithic ZIPs, debs, etc.
- [ ] Dynamic dependencies (code generation a la Swagger, etc.)
