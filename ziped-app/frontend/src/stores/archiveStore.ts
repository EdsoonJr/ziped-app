import { writable, derived } from 'svelte/store';
import type { models } from '../../wailsjs/go/models';

export const currentArchivePath = writable<string | null>(null);
export const currentArchiveFiles = writable<models.FileInfo[]>([]);
export const selectedFiles = writable<string[]>([]);
export const isProcessing = writable<boolean>(false);
export const progressPercent = writable<number>(-1);
export const currentProgressFile = writable<string>("");
export const progressStatus = writable<string>("");

export const isReadOnly = derived(currentArchivePath, $path => {
    if (!$path) return false;
    const ext = $path.split('.').pop()?.toLowerCase();
    return ext === 'rar' || ext === '7z';
});
