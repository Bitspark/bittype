import type { External as E } from "./external";
import type * as Other from "./external";

export interface Person { name: string; readonly age?: number; readonly symbol: unique symbol; }
export interface Page<T extends Person = Person> { items: T[]; next?: string; }
export interface Reader<T> {
  read(key: string, ...rest: number[]): T;
  <U>(value: U): U;
  readonly [key: string]: unknown;
}
export type Applied = Page<Person>;
type Qualified = Other.External;
type Literal = "a\\b" | 'escaped\'quote' | 12.5 | -2 | 42n | true | false | null;
type Tuple = [head: Person, tail?: string];
type Rest = [...values: string[]];
type Unnamed = [string, number?];
type Fn = <T>(value: T) => T;
type Functions = ((value: string) => number)[];
type Keys = keyof Person;
type Immutable = readonly string[];
type Access = Page<Person>["items"];
type Mix = Person & { "unusual-key"?: E; optional?(): void };
type Primitives = any | unknown | never | void | undefined | boolean | number | bigint | string | symbol | object;
interface Derived extends Page<Person> { extra: string; }
