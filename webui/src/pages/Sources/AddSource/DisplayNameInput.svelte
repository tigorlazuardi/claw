<script lang="ts">
  interface Props {
    value: string;
  }

  let { value = $bindable("") }: Props = $props();

  const allOk = $derived(!!value);
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span class={{ "text-success": allOk }}>
      Display Name <span class="text-error">*</span>
    </span>
  </legend>
  <input
    type="text"
    class={{
      "input w-full": true,
      "input-success": allOk,
      "text-success": allOk,
    }}
    bind:value
    placeholder="Friendly name to be shown in the UI for this source"
    required
  />
  <p
    class={{
      label: true,
      "text-success": allOk,
    }}
  >
    Required
  </p>
</fieldset>
