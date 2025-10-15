<script lang="ts">
  import { CreateDeviceSchema, CreateDeviceUiSchema } from "./schema";
  import { createForm, BasicForm } from "@sjsf/form";
  import { resolver } from "@sjsf/form/resolvers/basic";
  import { translation } from "@sjsf/form/translations/en";
  import { createFormValidator } from "@sjsf/ajv8-validator";
  import { theme } from "@sjsf/daisyui5-theme";
  import DialogModal from "#/components/DialogModal.svelte";

  interface Props {
    open: boolean;
  }

  let { open = $bindable(false) }: Props = $props();

  const validator = createFormValidator({ uiSchema: CreateDeviceUiSchema });

  const form = createForm({
    schema: CreateDeviceSchema,
    resolver,
    translation,
    validator,
    theme,
    uiSchema: CreateDeviceUiSchema,
    onSubmit: console.info,
  });
</script>

<DialogModal
  contentProps={{ interactOutsideBehavior: "ignore" }}
  class="w-[90vw] sm:w-[60vw]"
  bind:open
>
  {#snippet title()}
    <span class="font-bold">Add New Device</span>
  {/snippet}
  {#snippet description()}
    <BasicForm {form}></BasicForm>
  {/snippet}
</DialogModal>
