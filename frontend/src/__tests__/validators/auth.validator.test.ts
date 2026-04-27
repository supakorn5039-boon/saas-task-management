import { describe, expect, it } from "vitest";
import { loginSchema, registerSchema } from "@/validators/auth.validator";

describe("loginSchema", () => {
  it("accepts a valid email and password", () => {
    const result = loginSchema.safeParse({
      email: "alice@example.com",
      password: "password",
    });
    expect(result.success).toBe(true);
  });

  it("rejects malformed email", () => {
    const result = loginSchema.safeParse({
      email: "not-an-email",
      password: "password",
    });
    expect(result.success).toBe(false);
  });
});

describe("registerSchema", () => {
  it("accepts matching passwords ≥ 8 chars", () => {
    const result = registerSchema.safeParse({
      email: "bob@example.com",
      password: "longenough",
      confirmPassword: "longenough",
    });
    expect(result.success).toBe(true);
  });

  it("rejects short password", () => {
    const result = registerSchema.safeParse({
      email: "bob@example.com",
      password: "short",
      confirmPassword: "short",
    });
    expect(result.success).toBe(false);
  });

  it("rejects mismatched confirmation on the confirmPassword field", () => {
    const result = registerSchema.safeParse({
      email: "bob@example.com",
      password: "longenough",
      confirmPassword: "different",
    });
    expect(result.success).toBe(false);
    if (!result.success) {
      const path = result.error.issues[0]?.path;
      expect(path).toContain("confirmPassword");
    }
  });
});
