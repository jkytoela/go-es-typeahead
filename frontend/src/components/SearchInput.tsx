import * as React from 'react';
import type { ComboBoxProps } from '@react-types/combobox';
import { useComboBoxState, useSearchFieldState } from 'react-stately';
import { useComboBox, useButton, useSearchField } from 'react-aria';
import { SearchIcon, XIcon } from '@heroicons/react/solid';
import { ListBox } from './ListBox';
import { Popover } from './Popover';

export function SearchInput<T extends object>(props: ComboBoxProps<T>) {
  const state = useComboBoxState({ ...props });
  const inputRef = React.useRef(null);
  const listBoxRef = React.useRef(null);
  const popoverRef = React.useRef(null);

  const { inputProps, listBoxProps, labelProps } = useComboBox(
    {
      ...props,
      inputRef,
      listBoxRef,
      popoverRef
    },
    state
  );

  React.useEffect(() => {
    if ([...state.collection]?.length > 0) {
      state.open();
    } else {
      state.close();
    }
  }, [state]);

  // Get props for the clear button from useSearchField
  const searchProps = {
    label: props.label,
    value: state.inputValue,
    onChange: (v: string) => state.setInputValue(v)
  };

  const searchState = useSearchFieldState(searchProps);
  const { clearButtonProps } = useSearchField(searchProps, searchState, inputRef);
  const clearButtonRef = React.useRef(null);
  const { buttonProps } = useButton(clearButtonProps, clearButtonRef);

  return (
    <div className="inline-flex flex-col relative w-full max-w-2xl">
      <label
        {...labelProps}
        className="block text-sm font-medium text-gray-700 text-left"
      >
        {props.label}
      </label>
      <div
        className={`relative px-2 inline-flex flex-row items-center rounded-t-2xl overflow-hidden shadow-sm border ${
          state.isFocused ? 'border-gray-300' : 'border-gray-300'
        } ${
          state.isOpen ? 'rounded-b-none border-b-0' : 'rounded-b-2xl'
        }`}
      >
        <SearchIcon aria-hidden="true" className="w-5 h-5 text-gray-500" />
        <input
          {...inputProps}
          ref={inputRef}
          className="outline-none px-3 py-1 appearance-none w-full"
        />
        <button
          {...buttonProps}
          ref={clearButtonRef}
          style={{ visibility: state.inputValue !== '' ? 'visible' : 'hidden' }}
          className="cursor-default text-gray-500 hover:text-gray-600"
        >
          <XIcon aria-hidden="true" className="w-4 h-4" />
        </button>
      </div>
      {[...state?.collection].length ? (
        <Popover
          popoverRef={popoverRef}
          isOpen={state.isOpen}
          onClose={state.close}
        >
          <ListBox {...listBoxProps} listBoxRef={listBoxRef} state={state} />
        </Popover>
      ): null}
    </div>
  );
}
