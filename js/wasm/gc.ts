import JsBridgeType, {KEY_INTERNAL_ID} from "./JsBridgeType.js";
import {deleteObject, ObjectHandle, RuntimeSharedDataContext} from "./cmd.js";

const fr = new FinalizationRegistry<[RuntimeSharedDataContext, ObjectHandle]>(([ctx, handle]) => {
    deleteObject(ctx, handle)
})

export function registerAutoGc(ctx: RuntimeSharedDataContext, jbt: JsBridgeType) {
    fr.register(jbt, [ctx, jbt[KEY_INTERNAL_ID]])
}
