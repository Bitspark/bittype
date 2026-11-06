pub trait Family {
    type Item<T>;
}
pub struct VectorFamily {}
impl Family for VectorFamily {
    type Item<T> = Vec<T>;
}
pub type Bytes = <VectorFamily as Family>::Item<u8>;

pub trait Lending {
    type Item<'a>: 'a where Self: 'a;
}
pub struct Borrowed<T> { pub value: T }
impl<T> Lending for Borrowed<T> {
    type Item<'a> = &'a T where Self: 'a;
}
pub type Borrow<'a, T: 'a> = <Borrowed<T> as Lending>::Item<'a>;

pub trait Transform<T> where T: Clone {
    type Output<U: Clone>: Clone where U: Default;
}
impl<T: Clone> Transform<T> for Borrowed<T> where T: Default {
    type Output<U: Clone> = Vec<U> where U: Default;
}

pub trait View {
    type Target: ?Sized;
    type Borrow<'a, T: ?Sized + 'a>: ?Sized + 'a where Self: 'a;
}
impl View for VectorFamily {
    type Target = str;
    type Borrow<'a, T: ?Sized + 'a> = T where Self: 'a;
}

pub trait BorrowPair {
    type Pair<'a: 'b, 'b, T: 'a>: 'b;
}
impl BorrowPair for VectorFamily {
    type Pair<'a: 'b, 'b, T: 'a> = (&'a T, &'b T);
}

pub trait Raw { type r#type<T>; }
impl Raw for () { type r#type<T> = T; }

mod nested {
    pub trait Marker {}
    impl Marker for u8 {}
    impl<'a, T: ?Sized> Marker for &'a T {}
    impl<T> Marker for [T; 2] {}
}
