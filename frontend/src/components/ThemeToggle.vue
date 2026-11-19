<template>
    <button
        type="button"
        class="p-1.5 rounded-full hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors"
        :title="darkMode ? 'Light mode' : 'Dark mode'"
        :aria-label="darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
        @click="toggle"
    >
        <IconSun v-if="darkMode" class="size-4" />
        <IconMoon v-else class="size-4" />
    </button>
</template>

<script setup lang="ts">
import { ref } from "vue";
import IconMoon from "~icons/lucide/moon";
import IconSun from "~icons/lucide/sun";

// The initial theme is set in index.html before the app loads, to avoid a flash of the wrong theme
const darkMode = ref(document.documentElement.getAttribute("data-theme") === "dark");

const toggle = () => {
    darkMode.value = !darkMode.value;
    const theme = darkMode.value ? "dark" : "light";
    document.documentElement.setAttribute("data-theme", theme);
    try {
        localStorage.setItem("theme", theme);
    } catch {
        // Storage unavailable (e.g. private mode) - the theme just won't persist
    }
};
</script>
