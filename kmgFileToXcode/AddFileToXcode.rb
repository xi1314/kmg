require "xcodeproj"
def AddFileToProject(file_path,path)
  unless Pathname.new(path).exist?
    raise "[Xcodeproj] Unable to open `#{path}` because it doesn't exist."
  end
  project = Xcodeproj::Project.new(path, true)
  project.send(:initialize_from_file)
  target = project.targets.first
  group = project.main_group.groups.first
  group.set_source_tree('SOURCE_ROOT')
  file_ref = group.new_reference(file_path)
  target.add_file_references([file_ref])
  project.save
end
filepath = ARGV.first
projectpath = ARGV[1]
AddFileToProject(filepath,projectpath)
