use std::vec::Vec as Vector;
/* outside /* nested */ comment */
mod nested { pub(crate) type Alias = u8; }
pub struct Person { pub name: String }
pub struct Page<T> { pub items: Vec<T>, pub next: Option<String> }
pub type People = Page<Person>;
pub struct Pair<'a, T: 'a>(pub &'a T, pub *const T) where T: Sized;
pub struct Defaults<T = Person> { pub value: T }
pub enum Choice<T> { Empty, One(T), Named { value: T } }
pub trait Reader<'a, T>: Send where T: 'a {
    type Item: Clone;
    fn read<U>(&'a self, key: &str, value: U) -> T where U: Clone;
    fn own(self);
    fn edit(&mut self);
}
pub trait Lifetimes<'a: 'b, 'b, T: Clone + 'a> {}
type Forms<'a> = (&'a mut [u8], *mut u8, [u8; 3], (), unsafe extern "C" fn(i32) -> i32, &'a dyn Send);
type Diverge = fn() -> !;
type Absolute = ::std::string::String;
type NamedFn = fn(value: i32);
