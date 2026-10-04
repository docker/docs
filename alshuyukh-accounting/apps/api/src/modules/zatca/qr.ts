/**
 * ZATCA QR code: a base64 TLV (tag, length, value) list.
 *   1 seller name · 2 VAT number · 3 timestamp · 4 total with VAT · 5 VAT total
 *   (phase 2 adds) 6 invoice hash · 7 ECDSA signature · 8 ECDSA public key
 *   · 9 the CA's signature of the stamp certificate (simplified invoices only)
 */
export interface QrFields {
  sellerName: string; vatNumber: string; timestamp: string; total: string; vatTotal: string;
  invoiceHash?: string; signature?: string; publicKey?: Buffer; certificateSignature?: Buffer;
}

export function encodeQr(f: QrFields): string {
  const parts: [number, Buffer][] = [
    [1, Buffer.from(f.sellerName, 'utf8')], [2, Buffer.from(f.vatNumber)], [3, Buffer.from(f.timestamp)],
    [4, Buffer.from(f.total)], [5, Buffer.from(f.vatTotal)],
  ];
  if (f.invoiceHash) parts.push([6, Buffer.from(f.invoiceHash)]);
  if (f.signature) parts.push([7, Buffer.from(f.signature)]);
  if (f.publicKey) parts.push([8, f.publicKey]);
  if (f.certificateSignature) parts.push([9, f.certificateSignature]);
  for (const [tag, v] of parts) if (v.length > 255) throw new Error(`QR field ${tag} is longer than 255 bytes`);
  return Buffer.concat(parts.map(([tag, v]) => Buffer.concat([Buffer.from([tag, v.length]), v]))).toString('base64');
}

export function decodeQr(base64: string): Map<number, Buffer> {
  const buf = Buffer.from(base64, 'base64');
  const out = new Map<number, Buffer>();
  for (let i = 0; i < buf.length;) {
    const tag = buf[i]!;
    const len = buf[i + 1]!;
    out.set(tag, buf.subarray(i + 2, i + 2 + len));
    i += 2 + len;
  }
  return out;
}
