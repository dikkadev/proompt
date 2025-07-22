// Color utilities for generating variations and managing theme colors

export interface ColorHSL {
  h: number;
  s: number;
  l: number;
}

export function hexToHsl(hex: string): ColorHSL {
  // Remove # if present
  hex = hex.replace('#', '');
  
  // Parse r, g, b values
  const r = parseInt(hex.substr(0, 2), 16) / 255;
  const g = parseInt(hex.substr(2, 2), 16) / 255;
  const b = parseInt(hex.substr(4, 2), 16) / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0;
  let s = 0;
  const l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    l: Math.round(l * 100)
  };
}

export function hslToString(hsl: ColorHSL): string {
  return `${hsl.h} ${hsl.s}% ${hsl.l}%`;
}

export function generateColorVariations(baseColor: ColorHSL) {
  return {
    primary: hslToString(baseColor),
    primarySoft: hslToString({ ...baseColor, s: Math.max(baseColor.s * 0.3, 10), l: Math.min(baseColor.l + 30, 95) }),
    primaryHover: hslToString({ ...baseColor, l: Math.max(baseColor.l - 5, 10) }),
    primaryMuted: hslToString({ ...baseColor, s: Math.max(baseColor.s * 0.2, 5), l: Math.min(baseColor.l + 40, 90) }),
  };
}

export function applyThemeColor(color: string) {
  const hsl = hexToHsl(color);
  const variations = generateColorVariations(hsl);
  
  const root = document.documentElement;
  
  // Apply color variations for light mode
  root.style.setProperty('--primary', variations.primary);
  root.style.setProperty('--primary-soft', variations.primarySoft);
  root.style.setProperty('--primary-hover', variations.primaryHover);
  root.style.setProperty('--ring', variations.primary);
  
  // Update variable states to use the accent color
  root.style.setProperty('--variable-provided', variations.primary);
  
  // Create dark mode variations (brighter and adjusted for dark backgrounds)
  const darkVariations = {
    primary: hslToString({ ...hsl, l: Math.min(hsl.l + 10, 85) }), // Brighter for dark mode
    primarySoft: hslToString({ ...hsl, s: Math.max(hsl.s * 0.8, 20), l: Math.max(hsl.l * 0.3, 15) }),
    primaryHover: hslToString({ ...hsl, l: Math.min(hsl.l + 15, 90) }),
    ring: hslToString({ ...hsl, l: Math.min(hsl.l + 10, 85) })
  };
  
  // Apply dark mode specific styles using CSS custom properties
  const styleId = 'proompt-theme-colors';
  let styleEl = document.getElementById(styleId) as HTMLStyleElement;
  
  if (!styleEl) {
    styleEl = document.createElement('style');
    styleEl.id = styleId;
    document.head.appendChild(styleEl);
  }
  
  styleEl.textContent = `
    .dark {
      --primary: ${darkVariations.primary};
      --primary-soft: ${darkVariations.primarySoft};
      --primary-hover: ${darkVariations.primaryHover};
      --ring: ${darkVariations.ring};
      --variable-provided: ${darkVariations.primary};
      --sidebar-primary: ${darkVariations.primary};
      --sidebar-ring: ${darkVariations.primary};
    }
  `;
  
  // Store in cookie
  document.cookie = `proompt-accent-color=${color}; path=/; max-age=${365 * 24 * 60 * 60}`; // 1 year
}

export function getStoredAccentColor(): string | null {
  const cookies = document.cookie.split(';');
  for (let cookie of cookies) {
    const [name, value] = cookie.trim().split('=');
    if (name === 'proompt-accent-color') {
      return value;
    }
  }
  return null;
}

export function getDefaultAccentColor(): string {
  return '#1e9fa2'; // Default teal color
}