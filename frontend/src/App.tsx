import React from 'react';
import { SearchInput } from './components/SearchInput';
import { Item } from 'react-stately';
import useDebounce from './hooks/useDebounce';

function App() {
  const [items, setItems] = React.useState<string[]>([]);
  const [value, setValue] = React.useState<string>('')
  const queryStr = useDebounce<string>(value, 100)

  React.useEffect(() => {
    const fetcher = async () => {
      try {
        const res = await fetch(`http://localhost:8080/search?query=${queryStr}`);
        const { data } = await res.json();
        setItems(data);
      } catch (error) {
        console.error('TODO: Add error handling...', error);
      }
    };

    if (queryStr.length > 0) {
      fetcher();
    } else {
      setItems([]);
    }
  }, [queryStr]);

  return <div className="flex items-center justify-center pt-10">
    <SearchInput label="Search" onInputChange={setValue}>
      {items?.map((item) => <Item key={item}>{item}</Item>)}
    </SearchInput>
  </div>;
}

export default App;
