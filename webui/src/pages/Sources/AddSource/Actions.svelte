<script lang="ts">
  type ButtonEvent = MouseEvent & {
    currentTarget: EventTarget & HTMLButtonElement;
  };

  interface Props {
    onclick: (event: ButtonEvent, immediate: boolean) => void;
    oncancel: (event: ButtonEvent) => void;
    disabled?: boolean;
  }

  export const { onclick, oncancel, disabled = false }: Props = $props();
</script>

<div class="modal-action">
  <button
    onclick={(evt) => onclick(evt, false)}
    type="button"
    class="btn btn-primary"
    {disabled}
  >
    Save
  </button>
  <button
    type="button"
    onclick={(evt) => onclick(evt, true)}
    class="btn btn-secondary"
    {disabled}
  >
    Save {"&"} Run
  </button>
  <button type="button" class="btn btn-error" onclick={oncancel}>Cancel</button>
</div>
