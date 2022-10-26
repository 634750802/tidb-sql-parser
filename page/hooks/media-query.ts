import {onActivated, onDeactivated, onMounted, onUnmounted, ref} from "vue";

export function mediaQuery(query: string) {
    const ml = matchMedia(query)
    const matches = ref(ml.matches)

    const handleChange = () => {
        matches.value = ml.matches
    }

    const addListener = () => {
        ml.addEventListener('change', handleChange)
    }

    const removeListener = () => {
        ml.removeEventListener('change', handleChange)
    }

    onMounted(addListener)
    onActivated(addListener)
    onDeactivated(removeListener)
    onUnmounted(removeListener)

    return matches
}

export function preferDarkScheme() {
    return mediaQuery("(prefers-color-scheme: dark)")
}