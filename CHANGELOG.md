# TBD
* Make tests declare a timeout and mark them as failed if they don't complete in that time
* Add the controller log level (as string) as an argument for constructing a TestSuiteRunner - this string will get passed as-is to the controller image using the special environment variable `LOG_LEVEL`

# 0.3.1
* explicitly specify service IDs in network configurations

# 0.3.0
* Stop the node network after the test controller runs
* Rename ServiceFactory(,Config) -> ServiceInitializer(,Core)
* Fix bug with not actually catching panics when running a test
* Fix a bug with the TestController not appropriately catching panics
* Log test result as soon as the test is finished
* Add some extra unit tests
* Implement controller log propagation
* Allow services to declare arbitrary file-based dependencies (necessary for staking)
