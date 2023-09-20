'use client'
import { Checkbox as ChakraCheckbox, FormLabel, Flex} from "@chakra-ui/react";

interface Props {
    label?: string
    name?: string,
    isChecked?:boolean
    onChange?: (event: React.FormEvent<HTMLInputElement>) => void
}

const Checkbox = ( { label, name, isChecked, onChange } : Props) => {
    return (
        <Flex>
            {label && <FormLabel>{label}</FormLabel>}
            <ChakraCheckbox
                name={name}
                isChecked={isChecked}
                onChange={onChange}
            />
        </Flex>
    );
};

export default Checkbox;