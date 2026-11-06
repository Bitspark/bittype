use std::{env, process::Command};

fn main() {
    let output = Command::new(env::var("RUSTC").expect("explicit Cargo compiler"))
        .arg("--version")
        .output()
        .expect("compiler version");
    let version = String::from_utf8(output.stdout).expect("compiler version UTF-8");
    assert!(output.status.success() && version.starts_with("rustc 1.90.0 ("));
    assert_eq!(env::var("HOST").unwrap(), "x86_64-unknown-linux-gnu");
    assert_eq!(env::var("TARGET").unwrap(), "x86_64-unknown-linux-gnu");
    println!("cargo:rustc-env=BITTYPE_RUSTC_VERSION={}", version.trim());
}
