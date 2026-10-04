/**
 * Minimal ASN.1 DER encoder and reader, enough to build a PKCS#10 CSR with
 * ZATCA's extensions and to read fields out of an X.509 certificate.
 */

function length(n: number): Buffer {
  if (n < 0x80) return Buffer.from([n]);
  const bytes: number[] = [];
  for (let v = n; v > 0; v = Math.floor(v / 256)) bytes.unshift(v & 0xff);
  return Buffer.from([0x80 | bytes.length, ...bytes]);
}

export const tlv = (tag: number, value: Buffer) => Buffer.concat([Buffer.from([tag]), length(value.length), value]);
export const seq = (...items: Buffer[]) => tlv(0x30, Buffer.concat(items));
export const set = (...items: Buffer[]) => tlv(0x31, Buffer.concat(items));
export const int = (n: number) => tlv(0x02, Buffer.from([n]));
export const utf8 = (s: string) => tlv(0x0c, Buffer.from(s, 'utf8'));
export const printable = (s: string) => tlv(0x13, Buffer.from(s, 'ascii'));
export const octets = (b: Buffer) => tlv(0x04, b);
export const bits = (b: Buffer) => tlv(0x03, Buffer.concat([Buffer.from([0]), b]));
/** Context-specific constructed tag [n]. */
export const ctx = (n: number, ...items: Buffer[]) => tlv(0xa0 | n, Buffer.concat(items));

export function oid(dotted: string): Buffer {
  const parts = dotted.split('.').map(Number);
  const out = [40 * parts[0]! + parts[1]!];
  for (const p of parts.slice(2)) {
    const stack: number[] = [p & 0x7f];
    for (let v = Math.floor(p / 128); v > 0; v = Math.floor(v / 128)) stack.unshift((v & 0x7f) | 0x80);
    out.push(...stack);
  }
  return tlv(0x06, Buffer.from(out));
}

export interface Node { tag: number; start: number; header: number; length: number; value: Buffer; raw: Buffer }

/** Reads one TLV node at `offset`. */
export function read(buf: Buffer, offset = 0): Node {
  const tag = buf[offset]!;
  let len = buf[offset + 1]!;
  let header = 2;
  if (len & 0x80) {
    const n = len & 0x7f;
    len = 0;
    for (let i = 0; i < n; i++) len = len * 256 + buf[offset + 2 + i]!;
    header += n;
  }
  return { tag, start: offset, header, length: len, value: buf.subarray(offset + header, offset + header + len), raw: buf.subarray(offset, offset + header + len) };
}

/** Children of a constructed node. */
export function children(node: Node): Node[] {
  const out: Node[] = [];
  for (let off = 0; off < node.value.length;) {
    const c = read(node.value, off);
    out.push(c);
    off += c.header + c.length;
  }
  return out;
}
