import subprocess
import time
import pandas as pd
import shutil
import os

def combine_sheets(file_list, output_file='combined_excel.xlsx'):
    # Dictionary to keep track of dataframes for each sheet name
    sheets_dict = {}

    for file_index, file_name in enumerate(file_list):
        xls = pd.ExcelFile(file_name)
        for sheet_name in xls.sheet_names:
            df = pd.read_excel(xls, sheet_name)
            df['Source'] = f'File_{file_index + 1}'  # Add a column to indicate the source file
            df = pd.concat([df.iloc[:, :1], df.iloc[:, 1:].reindex(sorted(df.columns[1:]), axis=1)], axis=1)
            # Append data to the respective sheet's dataframe in the dictionary
            if sheet_name not in sheets_dict:
                sheets_dict[sheet_name] = df
            else:
                sheets_dict[sheet_name] = pd.concat([sheets_dict[sheet_name], df], ignore_index=True)

    # Write each dataframe to its corresponding sheet in the combined Excel file
    with pd.ExcelWriter(output_file, engine='openpyxl') as writer:
        for sheet_name, dataframe in sheets_dict.items():
            dataframe.to_excel(writer, sheet_name=sheet_name, index=False)

# Path to the directory containing the GoLang script and statistics.xlsx
golang_script_directory = r"C:\Users\rohan\OneDrive - Imperial College London\Imperial Year 4\SOMAS\experiments\SOMAS2023"

games = 10  # Specify the number of iterations
excel_files = []
original_file_name = 'statistics.xlsx'

for i in range(games):
    # Run the GoLang script that generates statistics.xlsx
    subprocess.run(['go', 'run', 'main.go'], cwd=golang_script_directory)

    # Copy and rename the newly generated statistics.xlsx to a new name
    new_file_name = f'statistics{i+1}.xlsx'
    new_file_path = os.path.join(golang_script_directory, new_file_name)
    shutil.copyfile(os.path.join(golang_script_directory, original_file_name), new_file_path)

    excel_files.append(new_file_path)

    time.sleep(1)  # Wait a bit before the next iteration (if necessary)

# Combine the generated Excel files
combine_sheets(excel_files)
print("Excel sheets combined successfully.")