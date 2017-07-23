use std::path::PathBuf;

#[derive(Clone)]
pub struct Dependency {
    pub name: String,
    pub ver: String
}

pub struct Repository {
    pub name: String,
    root: PathBuf
}

impl Repository {

    pub fn new(name: String, path: PathBuf) -> Repository {
        Repository { name: name, root: path }
    }

    pub fn get_package(&self, dep: Dependency) -> Option<PathBuf> {

        let mut path = self.root.clone();
        path.push(dep.name);
        path.push(dep.ver);

        if path.exists() {
            Some(path)
        } else {
            None
        }

    }

}
