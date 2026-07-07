import { describe, expect, test } from "bun:test";
import { cell } from "../src/text";

describe("cell", () => {
  test("pads and truncates ASCII", () => {
    expect(cell("api", 6)).toBe("api   ");
    const got = cell("devspace", 6);
    expect(got).toBe("devs… ");
    expect(Bun.stringWidth(got)).toBe(6);
  });

  test("handles emoji width", () => {
    const got = cell("api 🚀", 8);
    expect(got).toContain("🚀");
    expect(Bun.stringWidth(got)).toBe(8);
  });

  test("truncates CJK by display width", () => {
    const got = cell("日本語プロジェクト", 10);
    expect(got.endsWith(" ")).toBe(true);
    expect(got).toContain("…");
    expect(Bun.stringWidth(got)).toBe(10);
  });
});
