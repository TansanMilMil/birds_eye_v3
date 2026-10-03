import { useCallback, useSyncExternalStore } from "react";
import { LocalStorageKey } from "../localStorage/localStorageKey";
import { News } from "../types/news";

const listeners = new Set<() => void>();
let cache: { raw: string | null; items: News[] } | null = null;

const readRaw = (): string | null => {
  try {
    return localStorage.getItem(LocalStorageKey.Bookmarks);
  } catch {
    return null;
  }
};

const isNews = (v: unknown): v is News =>
  typeof v === "object" &&
  v !== null &&
  typeof (v as News).articleUrl === "string";

const parse = (raw: string | null): News[] => {
  if (!raw) return [];
  try {
    const parsed: unknown = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed.filter(isNews) : [];
  } catch {
    return [];
  }
};

const getSnapshot = (): News[] => {
  const raw = readRaw();
  if (!cache || cache.raw !== raw) {
    cache = { raw, items: parse(raw) };
  }
  return cache.items;
};

const subscribe = (listener: () => void) => {
  listeners.add(listener);
  window.addEventListener("storage", listener);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", listener);
  };
};

const write = (items: News[]) => {
  try {
    localStorage.setItem(LocalStorageKey.Bookmarks, JSON.stringify(items));
  } catch {
    return;
  }
  listeners.forEach((l) => l());
};

export function useBookmarks() {
  const bookmarks = useSyncExternalStore(subscribe, getSnapshot);

  const isBookmarked = useCallback(
    (articleUrl: string) => bookmarks.some((n) => n.articleUrl === articleUrl),
    [bookmarks]
  );

  const toggleBookmark = useCallback((news: News) => {
    const current = getSnapshot();
    write(
      current.some((n) => n.articleUrl === news.articleUrl)
        ? current.filter((n) => n.articleUrl !== news.articleUrl)
        : [news, ...current]
    );
  }, []);

  return { bookmarks, isBookmarked, toggleBookmark };
}
