import { useAuthStore } from '../stores/authStore';
import type { Theme } from '../types/auth';

export const useTheme = () => {
  const { theme, setTheme, updateTheme } = useAuthStore();

  const toggleTheme = async () => {
    const newTheme: Theme = theme === 'light' ? 'dark' : 'light';
    await updateTheme(newTheme);
  };

  return {
    theme,
    setTheme,
    updateTheme,
    toggleTheme,
  };
};
