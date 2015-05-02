require 'rake/clean'

# special case: we have a bindata file which should be regenerated if source
# file changes.  this is needed during development only, so we handle that here
# versus in the build script.
BINDATA    = "commands/inits/bindata.go"
file BINDATA => :generate
CLEAN.include FileList.new("tmp/*")

# convenience bootstrap all for getting started
desc "bootstrap all gotool dependencies"
task :bootstrap do
  sh "script/bootstrap"
end

# runs the generate script, which will bootstrap anything it needs in script
desc "generates bindata files"
task :generate do
  sh "script/generate"
end

# the unix build script does not force `generate` as prereq, but the task here
# does since we want to always make sure to be up to date with any changes made
desc "builds the binary"
task :build => :generate do
  sh "script/build"
end

desc "builds & installs the binary to $GOPATH/bin"
task :install => :build do
  cp "bin/scmpuff", "#{ENV['GOPATH']}/bin/scmpuff"
end

desc "run unit tests"
task :test => BINDATA do
  sh "go test ./..."
end

desc "run integration tests"
task :features => :build do
  sh "cucumber -s --tags=~@wip"
end

task "features:wip" => :build do
  sh "cucumber -s --tags=@wip"
end

desc "package for distribution"
task :package do
  tagged_version = `script/version`.chomp()
  sh "goxc -pv='#{tagged_version}'"
end
CLOBBER.include "builds"

task :all => [:build, :test, :features]
task :default => :all
