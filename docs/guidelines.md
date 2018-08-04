# Testing Guidelines

Ensure that all go routines with an ```Expect``` must call ```GinkgoRecover```. Also goroutines in tests must use the co-go concurrency patterns.

```
co.ParBegin(
    func() {
        c = DoSomething()
    },
    func() {
        defer GinkgoRecover()

        Expect(c).To(BeTrue())
    })
```
                
_See the testutils/testutils\_test.go for a template that can be used to begin writing tests_

To run tests, type ```ginkgo``` within the package that is being tested.

It is also advised to run tests with the race detector turned on to catch race conditions early on. This can be done by typing ```ginkgo --race```.

You can run all tests of the project by typing ```ginkgo -r``` within the republic-go directory. Alternatively, you can also run ```travis.sh``` locally. For this, you will need to install ```covermerge```.

```sh
go get -v github.com/loongy/covermerge
./.travis.sh
```
