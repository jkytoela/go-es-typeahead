import React from 'react';
import { SearchInput } from './components/SearchInput';
import { Item } from 'react-stately';

function App() {
  const [items, setItems] = React.useState<string[]>([]);

  const onInputChange = (value: string) => {
    if (value === 'test') {
      setItems(['example', 'second example']);
    } else {
      setItems([]);
    }
  };

  return <div className="flex items-center justify-center pt-10">
    <SearchInput label="Search" onInputChange={onInputChange}>
      {items.map((item) => <Item key={item}>{item}</Item>)}
    </SearchInput>
  </div>;
}

export default App;
