export function cell(value: string, width: number): string {
  if (width <= 0) return "";
  let text = value;
  if (Bun.stringWidth(text) > width - 1) {
    const limit = Math.max(0, width - 2);
    let used = 0;
    text = "";
    for (const ch of value) {
      const next = Bun.stringWidth(ch);
      if (used + next > limit) break;
      text += ch;
      used += next;
    }
    text += "…";
  }
  return text + " ".repeat(Math.max(0, width - Bun.stringWidth(text)));
}
