// Projections retain the subject, optional trait and suffix separately.
pub trait Collection { type Item; }
pub type Item<C: Collection> = <C as Collection>::Item;
pub type Inferred<C: Collection> = <C>::Item;
pub type Nested<C: Collection> where <C as Collection>::Item: Collection = <<C as Collection>::Item as Collection>::Item;
pub type Qualified = <std::vec::IntoIter<u8> as ::std::iter::Iterator>::Item;
pub type Turbofish = <std::vec::IntoIter::<u8> as ::std::iter::Iterator>::Item;
pub type Bytes = dyn Iterator<Item = u8>;
pub struct Iter<I: Iterator<Item = u8>> { pub value: I }
pub trait Source { type Value; }
pub trait Converted: Source + Iterator<Item = <Self as Source>::Value> {}
pub trait Relation<'a, T> { type Left; type Right; }
pub struct Bound<'a, T: 'a, R: Relation<'a, T, Left = T, Right = &'a T>> {
    pub value: R,
    pub marker: &'a T,
}
pub trait Raw { type r#type; }
pub type RawItem<T: Raw> = <T as Raw>::r#type;
pub trait RawBound: Raw<r#type = u8> {}
