require 'rake/clean'

# special case: we have a bindata file which should be regenerated if source
# file changes.  this is needed during development only, so we handle that here
# versus in the build script.
BINDATA    = "commands/inits/bindata.go"
SH_SCRIPTS = FileList.new("commands/inits/data/*.sh")

# the bindata file is defined as being dependent on all shell scripts in data/
# if any shell scripts change, clean intermediary files then regenerate bindata
file BINDATA => [*SH_SCRIPTS] do
  FileUtils.rm(BINDATA, :verbose => true) if File.exists? BINDATA
  sh "go generate ./commands/inits"
end

CLEAN.include(BINDATA) if File.exists? BINDATA
CLEAN.include FileList.new("tmp/*")

desc "bootstrap gotool dependencies"
task :bootstrap do
  sh "script/bootstrap"
end

desc "builds the binary"
task :build => BINDATA do
  sh "script/build"
end

desc "builds & installs the binary"
task :install => :build do
  sh "go install"
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

# quick and dirty package script, will have to replace with something much more
# robust and also cross platform down the line, but for sending to friends this
# can work for now...
desc "package for distribution"
task :package => [:build] do
  DEST = "builds/scmpuff-osx"
  tagged_version = `git describe --tags`.chomp()
  mkdir_p DEST
  cp "bin/scmpuff",   DEST
  cp "README.md",     DEST
  cp "CHANGELOG.md",  DEST
  cp "INSTALL.txt",   DEST
  sh "tar -C builds -cz -f builds/scmpuff_osx_amd64_#{tagged_version}.tgz scmpuff-osx"
end
CLOBBER.include "builds"

task :all => [:build, :test, :features]
task :default => :all
