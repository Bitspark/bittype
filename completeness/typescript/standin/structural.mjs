// PROTOTYPE (bittype#196). The run-time factories of the TypeScript stand-in
// for bittype's structural contributions. A plan entry names one of these by
// module and export, through the facet that crossed into the plan, and gives
// it the family's payload. They stand in for the service library's package
// `structural-adapters`, which does not exist. One seat wrote both sides, so
// nothing here shows independent authorship (row 31).
//
// A factory states the payload schemas it reads and the contribution
// references it implements. make receives the decoded payload, the node's
// links with what was built for each dependency, and the canonical name that
// the source snapshot gives a child.

export class Refusal extends Error {
  constructor(kind, message) {
    super(message);
    this.kind = kind;
  }
}

const schema = (name) => [{ identity: "bittype.payload/horos-to-typescript/" + name, version: "1" }];
const reference = (meaning) => ["bittype.structural/" + meaning + "/1#typescript/1"];
const own = (object, key) => Object.prototype.hasOwnProperty.call(object, key);
const plain = (v) => v !== null && typeof v === "object" && !Array.isArray(v);

export const text = {
  schemas: schema("text"),
  implements: reference("text"),
  make() {
    return {
      decode(v) {
        if (typeof v !== "string") throw new Refusal("structural mismatch", "text is expected");
        return v;
      },
      encode(n) {
        if (typeof n !== "string") throw new Refusal("invalid native value", "a string is expected, got " + (n === null ? "null" : typeof n));
        return n;
      },
    };
  },
};

export const list = {
  schemas: schema("list"),
  implements: reference("list"),
  make({ links }) {
    const element = links[0].adapter;
    return {
      decode(v) {
        if (!Array.isArray(v)) throw new Refusal("structural mismatch", "a list is expected");
        return v.map((item) => element.decode(item));
      },
      encode(n) {
        if (!Array.isArray(n)) throw new Refusal("invalid native value", "an array is expected");
        return n.map((item) => element.encode(item));
      },
    };
  },
};

// The adapter of a declaration is the adapter of its body. The recorded
// decision (nominal or alias, the parameter names) concerns the type only.
export const declaration = {
  schemas: schema("declarations"),
  implements: reference("declaration-policy"),
  make({ links }) {
    return links[0].adapter;
  },
};

// The record factory. It lays out the members that the record's payload
// names. For a field whose presence the record delegated it invokes the
// field-state protocol of what was selected for the field position: it never
// reads that contribution's payload. Only the record can omit a field, so it
// hands the protocol the container and the member's key, which the record
// stated in its request.
export const record = {
  schemas: schema("record"),
  implements: reference("closed-record"),
  make({ payload, links, canonical }) {
    let laid = 0;
    const fields = links.map((link, i) =>
      link.kind === "field"
        ? { name: canonical(i), member: link.stated, state: link.state }
        : { name: canonical(i), member: payload.members[laid++], adapter: link.adapter },
    );
    return {
      decode(v) {
        if (!plain(v)) throw new Refusal("structural mismatch", "a record is expected");
        for (const key of Object.keys(v)) {
          if (!fields.some((f) => f.name === key)) throw new Refusal("structural mismatch", `the record has no field "${key}"`);
        }
        const out = {};
        for (const f of fields) {
          if (f.state) {
            f.state.write(out, f.member, own(v, f.name), v[f.name]);
          } else {
            if (!own(v, f.name)) throw new Refusal("structural mismatch", `the required field "${f.name}" is absent`);
            out[f.member] = f.adapter.decode(v[f.name]);
          }
        }
        return out;
      },
      encode(n) {
        if (!plain(n)) throw new Refusal("invalid native value", "an object is expected");
        for (const key of Object.keys(n)) {
          if (!fields.some((f) => f.member === key)) throw new Refusal("invalid native value", `the object has no member "${key}"`);
        }
        const out = {};
        for (const f of fields) {
          if (f.state) {
            const state = f.state.read(n, f.member);
            if (state.present) out[f.name] = state.value;
          } else {
            if (!own(n, f.member)) throw new Refusal("invalid native value", `the required member "${f.member}" is absent`);
            out[f.name] = f.adapter.encode(n[f.member]);
          }
        }
        return out;
      },
    };
  },
};

// The field-state protocol of the optional contribution on the TypeScript
// side. It has no standalone value adapter (row 14).
//
// What it needs: the container object and the member's key, from the record,
// and the adapter selected for the operand. It needs the container because a
// TypeScript optional property has no value of its own when it is absent:
// absence is a fact about the object.
//
// The four states of a field (row 35, O14): absent; present with a value;
// present with null, which the operand's adapter refuses because no delivered
// family admits null; present with undefined, which is a present field whose
// value is no value of the field's type, and is refused here.
export const optionalField = {
  schemas: schema("optional"),
  implements: reference("optional-field"),
  make({ links }) {
    const operand = links[0].adapter;
    return {
      // write constructs or omits the native field from the field's state.
      write(container, member, present, value) {
        if (present) container[member] = operand.decode(value);
      },
      // read tests presence and obtains the present value.
      read(container, member) {
        if (!own(container, member)) return { present: false };
        const native = container[member];
        if (native === undefined) {
          if (process.env.BITTYPE_STANDIN_MUTANT === "undefined-as-absent") return { present: false };
          throw new Refusal("invalid native value", `the member "${member}" is present and undefined: a present undefined is no value of the field's type`);
        }
        return { present: true, value: operand.encode(native) };
      },
    };
  },
};
