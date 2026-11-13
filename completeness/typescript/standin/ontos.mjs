// PROTOTYPE (bittype#196). The part of ontos that the stand-in needs to read
// a frozen payload: the canonical codec and the data embeddings of text,
// booleans, lists and maps. A payload is an ontos value whose schema is a
// horos type; the stand-in reads it as data and validates nothing.

const TAG_ATOM = 0x00;
const TAG_TUPLE = 0x01;

export function fromHex(hex) {
  if (hex.length % 2 !== 0) throw new Error("ontos: odd hex length");
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0; i < bytes.length; i++) bytes[i] = parseInt(hex.slice(2 * i, 2 * i + 2), 16);
  return bytes;
}

export function decode(bytes) {
  let i = 0;
  const uvarint = () => {
    let n = 0;
    let shift = 0;
    for (;;) {
      if (i >= bytes.length) throw new Error("ontos: truncated");
      const b = bytes[i++];
      n += (b & 0x7f) * 2 ** shift;
      if (!(b & 0x80)) return n;
      shift += 7;
    }
  };
  const value = () => {
    const tag = bytes[i++];
    const n = uvarint();
    if (tag === TAG_ATOM) {
      const atom = bytes.subarray(i, i + n);
      i += n;
      return { atom };
    }
    if (tag === TAG_TUPLE) {
      const tuple = [];
      for (let k = 0; k < n; k++) tuple.push(value());
      return { tuple };
    }
    throw new Error("ontos: unknown tag " + tag);
  };
  const v = value();
  if (i !== bytes.length) throw new Error("ontos: trailing bytes");
  return v;
}

const text = (atom) => new TextDecoder("utf-8", { fatal: true }).decode(atom);

function labeled(v, label) {
  return v.tuple && v.tuple.length >= 1 && v.tuple[0].atom && text(v.tuple[0].atom) === label;
}

// data reads an ontos value of the embeddings the delivered families use into
// a JavaScript value: text, boolean, list, and a map with text keys.
export function data(v) {
  if (labeled(v, "utf8-text") && v.tuple.length === 2) return text(v.tuple[1].atom);
  if (labeled(v, "bool") && v.tuple.length === 2) return v.tuple[1].atom[0] === 1;
  if (labeled(v, "list")) return v.tuple.slice(1).map(data);
  if (labeled(v, "map")) {
    const out = {};
    for (const entry of v.tuple.slice(1)) out[data(entry.tuple[0])] = data(entry.tuple[1]);
    return out;
  }
  throw new Error("ontos: a value outside the embeddings the stand-in reads");
}
