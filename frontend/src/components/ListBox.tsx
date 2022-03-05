/* eslint-disable @typescript-eslint/no-use-before-define */
import * as React from 'react';
import type { AriaListBoxOptions } from '@react-aria/listbox';
import type { ListState } from 'react-stately';
import type { Node } from '@react-types/shared';
import { useListBox, useOption } from 'react-aria';

interface ListBoxProps extends AriaListBoxOptions<unknown> {
  listBoxRef?: React.RefObject<HTMLUListElement>;
  state: ListState<unknown>;
}

interface OptionProps {
  item: Node<unknown>;
  state: ListState<unknown>;
}

export function ListBox(props: ListBoxProps) {
  const ref = React.useRef<HTMLUListElement>(null);
  const { listBoxRef = ref, state } = props;
  const { listBoxProps } = useListBox(props, state, listBoxRef);

  return (
    <ul
      {...listBoxProps}
      ref={listBoxRef}
      className="max-h-72 overflow-auto outline-none pt-1 pb-6"
    >
      {[...state.collection].map((item) =>
        <Option key={item.key} item={item} state={state} />
      )}
    </ul>
  );
}

function Option({ item, state }: OptionProps) {
  const ref = React.useRef<HTMLLIElement>(null);
  const { optionProps, isSelected, isFocused } = useOption(
    {
      key: item.key
    },
    state,
    ref,
  );

  return (
    <li
      {...optionProps}
      ref={ref}
      className={`py-2 px-2 text-sm outline-none cursor-default flex items-center justify-between text-gray-900 ${
        isFocused ? 'bg-gray-200' : ''
      } ${isSelected ? 'font-bold' : ''}`}
    >
      {item.rendered}
    </li>
  );
}
