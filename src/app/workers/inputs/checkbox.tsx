'use client'
import { Checkbox as ChakraCheckbox, FormLabel, Flex} from "@chakra-ui/react";

interface Props {
    label?: string
    name?: string,
    isChecked?:boolean,
    disabled?:boolean,
    onChange?: (event:React.FormEvent<any>) => void
}

const Checkbox = ( { label, name, isChecked, disabled, onChange } : Props) => {
    return (
        <Flex>
            {label && <FormLabel>{label}</FormLabel>}
            <ChakraCheckbox
                name={name}
                isChecked={isChecked}
                onChange={onChange}
                disabled={disabled}
            />
        </Flex>
    );
};

export default Checkbox;