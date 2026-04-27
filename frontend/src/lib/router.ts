// Tiny shim that lets non-React modules (the axios interceptor) navigate via
// the TanStack router without a full page reload. main.tsx calls setRouter()
// during boot; anything that runs before that falls back to a hard redirect.
type RouterLike = { navigate: (opts: { to: string }) => unknown };

let router: RouterLike | null = null;

export function setRouter(r: RouterLike) {
  router = r;
}

export function navigateTo(to: string) {
  if (router) {
    router.navigate({ to });
  } else if (typeof window !== "undefined") {
    window.location.href = to;
  }
}
