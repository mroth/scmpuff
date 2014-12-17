require 'rake/clean'

# special case: we have a bindata file which should be regenerated if source
# file changes.  this is needed during development only, so we handle that here
# versus in the build script.
BINDATA    = "commands/inits/bindata.go"
SH_SCRIPTS = FileList.new("commands/inits/data/*.sh")

file BINDATA => [*SH_SCRIPTS] do
  FileUtils.rm(BINDATA, :verbose => true) if File.exists? BINDATA
  sh "go generate ./commands/inits"
end

# the bindata file is also considered an intermediary and can be cleaned up
CLEAN.include("commands/inits/bindata.go")


desc "bootstrap gotool dependencies"
task :bootstrap do
  sh "script/bootstrap"
end

desc "builds the binary"
task :build => BINDATA do
  sh "script/build"
end

desc "run unit tests"
task :test do
  sh "go test ./..."
end

desc "run integration tests"
task :features => :build do
  sh "cucumber -s --tags=~@wip"
end

task :all => [:build, :test, :features]
task :default => :all
