import { createContext, useContext } from 'react';
import de from './de';
import en from './en';

export type Locale = 'de' | 'en';
export type Translations = typeof de;

const translations: Record<Locale, Translations> = { de, en };

export function getTranslations(locale: Locale): Translations {
  return translations[locale] || translations.de;
}

export interface I18nContextType {
  locale: Locale;
  t: Translations;
  setLocale: (locale: Locale) => void;
}

export const I18nContext = createContext<I18nContextType>({
  locale: 'de',
  t: de,
  setLocale: () => {},
});

export function useI18n(): I18nContextType {
  return useContext(I18nContext);
}

export { de, en };
