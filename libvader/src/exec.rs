use std::path::{Path, PathBuf};
use std::env;
use std::fs;

use dep::*;

#[allow(dead_code)]
pub struct Execution {
    script: PathBuf,
    runtime: Runtime,
    deps: Vec<Box<Dependency>>
}

#[allow(dead_code)]
pub struct Runtime {
    name: String,
    path: PathBuf
}

impl Execution {

}

impl Runtime {

    pub fn get_system_runtimes() -> Result<Vec<Runtime>, &'static str> {

        match env::var("PATH") {
            Ok(ev_path) => {

                let dirs = ev_path.split(":"); // FIXME Windows is stupid and I don't remember how it works.
                let mut rts = Vec::new();

                for d in dirs {

                    let dp = PathBuf::from(Path::new(d));

                    // FIXME NOT FUNCTIONAL ENOUGH.
                    match fs::read_dir(dp.clone()) {
                        Ok(rd) => for e in rd {
                            match e {
                                Ok(ee) => match ee.file_name().into_string() {
                                    Ok(name) => if name.starts_with("python") {
                                        rts.push(Runtime { name: name, path: dp.clone() });
                                    },
                                    Err(_) => {}
                                },
                                Err(_) => {}
                            }
                        },
                        Err(_) => {}
                    }

                }

                Ok(rts)

            },
            Err(_) => Err("path not set")
        }

    }

}
