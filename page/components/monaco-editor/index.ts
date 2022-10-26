import TsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker.js?worker'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker.js?worker'

window.MonacoEnvironment = {
    getWorker: function (workerId, label) {
        switch (label) {
            case 'typescript':
            case 'javascript':
                return new TsWorker();
            default:
                return new EditorWorker();
        }
    }
};

export { default } from './MonacoEditor.vue';
