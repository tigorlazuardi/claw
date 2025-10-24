import { NSFWMode } from "#/gen/claw/v1/nsfw_pb";
import Type, { type Static } from "typebox";
import type { UiSchema } from "@sjsf/form";
import "@sjsf/form/fields/extra-fields/enum-include";
import SelectNSFWMode from "./SelectNSFWMode.svelte";

export const CreateDeviceSchema = Type.Object({
  slug: Type.String({
    minLength: 1,
    title: "Slug",
    description:
      "Unique identifier for the device. Must be url and filesystem friendly",
  }),
  name: Type.String({
    minLength: 1,
    title: "Name",
    description: "Human friendly name",
  }),
  height: Type.Integer({
    minimum: 1,
    title: "Height",
    description: "Device height in pixels",
  }),
  width: Type.Integer({ minimum: 1 }),
  aspectRatioDifference: Type.Number({ minimum: 0, default: 0.2 }),
  filenameTemplate: Type.Optional(Type.String()),
  imageMinHeight: Type.Optional(Type.Integer({ minimum: 0 })),
  imageMinWidth: Type.Optional(Type.Integer({ minimum: 0 })),
  imageMaxHeight: Type.Optional(Type.Integer({ minimum: 0 })),
  imageMaxWidth: Type.Optional(Type.Integer({ minimum: 0 })),
  imageMinFilesize: Type.Optional(Type.Integer({ minimum: 0 })),
  imageMaxFilesize: Type.Optional(Type.Integer({ minimum: 0 })),
  nsfw: Type.Enum(Object.keys(NSFWMode), {
    title: "NSFW Mode",
    description: "NSFW handling mode for this device",
  }),
  isDisabled: Type.Optional(
    Type.Boolean({
      default: false,
      title: "Disable Device",
      description: "If true, the device will not be available for use",
    }),
  ),
  sources: Type.Array(Type.Integer()),
});

export const CreateDeviceUiSchema: UiSchema = {
  nsfw: {
    "ui:components": {
      // TODO: Move to use custom Radio using bits-ui
      stringField: SelectNSFWMode,
    },
    "ui:options": {
      text: {
        autocomplete: "name",
      },
    },
  },
};

export type CreateDeviceSchemaType = Static<typeof CreateDeviceSchema>;
