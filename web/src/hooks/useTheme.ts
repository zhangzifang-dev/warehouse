import { message } from 'antd';
import { useAuthStore } from '../stores/authStore';
import type { Theme } from '../types/auth';

export const useTheme = () => {
  const { theme, setTheme, updateTheme } = useAuthStore();

  const toggleTheme = async () => {
    const newTheme: Theme = theme === 'light' ? 'dark' : 'light';
    try {
      await updateTheme(newTheme);
      message.success(`已切换到${newTheme === 'dark' ? '深色' : '浅色'}模式`);
    } catch (error) {
      console.error('Failed to update theme:', error);
      message.error('切换主题失败，请重试');
    }
  };

  return {
    theme,
    setTheme,
    updateTheme,
    toggleTheme,
  };
};
